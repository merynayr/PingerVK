package ping

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"

	"github.com/merynayr/PingerVK/pinger/internal/model"
	"github.com/merynayr/PingerVK/pkg/logger"
)

// Функция пингования IP-адресов контейнеров
func ping(ip string) (time.Duration, error) {
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return 0, fmt.Errorf("ошибка создания ICMP соединения: %v", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Ошибка при закрытии ICMP соединения: %v", err)
		}
	}()

	dst, err := net.ResolveIPAddr("ip4", ip)
	if err != nil {
		return 0, fmt.Errorf("не удалось разрешить IP %s: %v", ip, err)
	}

	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{ID: 1, Seq: 1, Data: []byte("hello")},
	}

	data, err := msg.Marshal(nil)
	if err != nil {
		return 0, fmt.Errorf("ошибка маршалинга ICMP сообщения: %v", err)
	}

	start := time.Now()
	_, err = conn.WriteTo(data, dst)
	if err != nil {
		return 0, fmt.Errorf("ошибка отправки ICMP запроса: %v", err)
	}

	reply := make([]byte, 1500)
	err = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	if err != nil {
		return 0, fmt.Errorf("ошибка установки тайм-аута: %v", err)
	}

	n, _, err := conn.ReadFrom(reply)
	if err != nil {
		return 0, fmt.Errorf("ответ не получен: %v", err)
	}

	duration := time.Since(start)

	resp, err := icmp.ParseMessage(1, reply[:n])
	if err != nil {
		return 0, fmt.Errorf("ошибка парсинга ICMP ответа: %v", err)
	}

	if resp.Type == ipv4.ICMPTypeEchoReply {
		return duration, nil
	}

	return 0, fmt.Errorf("неожиданный ICMP тип ответа: %v", resp.Type)
}

// Функция получения IP-адресов контейнеров
func getContainerIPs() []model.ContainerInfo {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logger.Error("Не удалось создать клиент Docker: ", err)
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		logger.Error("Не удалось получить список контейнеров: ", err)
	}

	var containerInfos []model.ContainerInfo
	for _, ctr := range containers {
		ctrInfo, err := cli.ContainerInspect(context.Background(), ctr.ID)
		if err != nil {
			log.Printf("Не удалось получить информацию о контейнере %s: %v", ctr.ID, err)
			continue
		}

		isRunning := ctrInfo.State.Running
		var ip string
		if isRunning {
			for _, netSettings := range ctrInfo.NetworkSettings.Networks {
				ip = netSettings.IPAddress
				break
			}
		}

		containerInfos = append(containerInfos, model.ContainerInfo{
			ID:        ctr.ID,
			Name:      ctrInfo.Name[1:],
			IPAddress: ip,
			Status:    isRunning,
		})
	}
	return containerInfos
}

// MonitorContainers Главная функция мониторинга контейнеров
func MonitorContainers(addres string) {
	containers := getContainerIPs()
	for _, container := range containers {
		if container.Name == "pinger" {
			continue
		}
		responseTime, err := ping(container.IPAddress)

		if err != nil {
			container.Status = false
		} else {
			container.Status = true
		}
		container.ResponseTime = responseTime
		if container.Status {
			container.LastSuccess = time.Now().UTC()
		}
		err = sendContainerStatus(container, addres)
		if err != nil {
			logger.Error("Ошибка отправки данных контейнера %s: %v", container.ID, err)
		}
	}
}

// Функция отправки POST-запроса с данными контейнера
func sendContainerStatus(container model.ContainerInfo, address string) error {
	url := fmt.Sprintf("http://%s/ping", address)

	data, err := json.Marshal(container)
	if err != nil {
		return fmt.Errorf("не удалось сериализовать данные контейнера: %v", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("ошибка при создании HTTP-запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка при отправке POST-запроса: %v", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("Ошибка при закрытии тела ответа: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неудачный статус ответа от сервера: %v", resp.Status)
	}

	return nil
}

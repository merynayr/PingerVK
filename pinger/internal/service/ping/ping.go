package ping

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/IBM/sarama"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"

	"github.com/merynayr/PingerVK/pinger/internal/model"
	"github.com/merynayr/PingerVK/pkg/logger"
)

// ping функция пингования IP-адресов контейнеров
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

// getContainerIPs функция получения IP-адресов контейнеров
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
func MonitorContainers() []model.ContainerInfo {
	containers := getContainerIPs()
	for i := range containers {
		if containers[i].Name == "pinger" {
			containers[i].Status = true
			continue
		}
		responseTime, err := ping(containers[i].IPAddress)

		if err != nil {
			containers[i].Status = false
		} else {
			containers[i].Status = true
		}
		containers[i].ResponseTime = responseTime
		if containers[i].Status {
			containers[i].LastSuccess = time.Now().UTC()
		}
	}
	return containers
}

// SendContainer функция отправки данных контейнера на backend
func (s *srv) SendContainer(pingTopicName string) error {
	containers := MonitorContainers()

	for _, container := range containers {
		data, err := json.Marshal(container)
		if err != nil {
			return fmt.Errorf("не удалось сериализовать данные контейнера: %v", err)
		}

		msg := &sarama.ProducerMessage{
			Topic:     pingTopicName,
			Partition: int32(0),
			Value:     sarama.StringEncoder(data),
		}

		res := s.kafkaProducer.SendMessage(msg)
		if res.Err != nil {
			logger.Error("failed to send message in Kafka", logger.With("error", err))
		}

		logger.Debug("message sent in Kafka", logger.With("container ip", container.IPAddress))
	}

	return nil
}

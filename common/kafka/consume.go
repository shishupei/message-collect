package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"strings"
)

type ConsumeConfig struct {
	Topic    string `yaml:"topic"`
	MqCert   string `yaml:"mqcert"`
	Address  string `yaml:"address"`
	Group    string `yaml:"group"`
	Offset   int64  `yaml:"offset"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
}

func ConsumeGroup(cfg ConsumeConfig, handler sarama.ConsumerGroupHandler) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = cfg.Offset
	config.Consumer.Return.Errors = true
	// 开始连接kafka服务器
	group, err := sarama.NewConsumerGroup(strings.Split(cfg.Address, ","), cfg.Group, config)

	if err != nil {
		fmt.Println("connect kafka failed; err", err)
		return
	}
	// 检查错误
	go func() {
		for err := range group.Errors() {
			fmt.Println("group errors : ", err)
		}
	}()

	ctx := context.Background()
	fmt.Println("start get msg " + cfg.Topic)
	// for 是应对 consumer rebalance

	// 需要监听的主题
	topics := []string{cfg.Topic}
	// 启动kafka消费组模式，消费的逻辑在上面的 ConsumeClaim 这个方法里
	err = group.Consume(ctx, topics, handler)
	if err != nil {
		fmt.Println("consume failed; err : ", err)
		return
	}
}

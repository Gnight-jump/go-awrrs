// Copyright (c) 2021 G_night

/*
	负载均衡 - 平滑的加权轮询算法
*/
package awwrs

import (
	"errors"
	"sync"
)

type WrrSlice struct {
	curAddr int
	nodes   []*WeightNode // 存储服务节点
	lock    sync.Mutex    // 锁，线程安全
}

type WeightNode struct {
	weight          int    // 初始化的节点权重
	currentWeight   int    // 节点当前权重
	effectiveWeight int    // 有效权重，用于判断健康状态，调用失败会-1，成功会+1，不会超过weight
	addr            string // 服务器的地址
}

/**
 * @Description：添加服务
 */
func (r *WrrSlice) Add(addr string, weight int) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	node := &WeightNode{
		weight:          weight,
		effectiveWeight: weight, // 初始化時有效权重 = 配置权重值
		currentWeight:   weight, // 初始化時当前权重 = 配置权重值
		addr:            addr,
	}
	r.nodes = append(r.nodes, node)
	return nil
}

/**
 * @Description：轮询获取服务
 */
func (r *WrrSlice) Next() (string, error) {
	// 锁，线程安全
	r.lock.Lock()
	defer r.lock.Unlock()

	// 服务数目为0
	if len(r.nodes) == 0 {
		return "", errors.New("[LOG_AWRRS] Inexistence service")
	}

	totalWeight := 0              // 总权重
	var maxWeightNode *WeightNode // 最大权重节点

	// 轮询获得最大权重
	for key, node := range r.nodes {
		totalWeight += node.effectiveWeight
		// 计算currentWeight
		node.currentWeight += node.effectiveWeight
		// 寻找权重最大的
		if maxWeightNode == nil || maxWeightNode.currentWeight < node.currentWeight {
			maxWeightNode = node
			r.curAddr = key
		}
	}

	if maxWeightNode != nil {
		// 更新选中节点的currentWeight
		maxWeightNode.currentWeight -= totalWeight
		// 返回addr
		return maxWeightNode.addr, nil
	}

	return "", errors.New("[LOG_AWRRS] The service node was not found")
}

package message

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/twotwo/go-blueprint/pkg/errors"
)

// CreateMessage 处理创建消息的请求
func CreateMessage(w http.ResponseWriter, r *http.Request) {
	var messageData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&messageData); err != nil {
		apiErr := errors.BadRequest("无效的请求体")
		errors.WriteJSON(w, apiErr)
		return
	}

	// 根据类型确定消息模型
	messageType, ok := messageData["type"].(string)
	if !ok {
		apiErr := errors.BadRequest("缺少消息类型或类型无效")
		errors.WriteJSON(w, apiErr)
		return
	}

	var responseMsg interface{}

	switch MessageType(messageType) {
	case Sms:
		var msg SMSMessage
		jsonData, err := json.Marshal(messageData)
		if err != nil {
			apiErr := errors.InternalServer("序列化请求数据失败")
			errors.WriteJSON(w, apiErr)
			return
		}

		if err := json.Unmarshal(jsonData, &msg); err != nil {
			apiErr := errors.BadRequest("无效的短信消息格式")
			errors.WriteJSON(w, apiErr)
			return
		}

		if msg.PhoneNumber == "" {
			apiErr := errors.BadRequest("短信消息必须包含手机号")
			errors.WriteJSON(w, apiErr)
			return
		}

		// 这里可以添加创建短信消息的业务逻辑
		msg.Id = getNextID() // 模拟ID生成
		msg.Status = Sending // 设置初始状态
		responseMsg = msg

	case Sitemessage:
		var msg SiteMessage
		jsonData, err := json.Marshal(messageData)
		if err != nil {
			apiErr := errors.InternalServer("序列化请求数据失败")
			errors.WriteJSON(w, apiErr)
			return
		}

		if err := json.Unmarshal(jsonData, &msg); err != nil {
			apiErr := errors.BadRequest("无效的站内消息格式")
			errors.WriteJSON(w, apiErr)
			return
		}

		if msg.UserId == 0 {
			apiErr := errors.BadRequest("站内消息必须包含用户ID")
			errors.WriteJSON(w, apiErr)
			return
		}

		// 这里可以添加创建站内消息的业务逻辑
		msg.Id = getNextID() // 模拟ID生成
		msg.Status = Sending // 设置初始状态
		responseMsg = msg

	case Broadcast:
		var msg BroadcastMessage
		jsonData, err := json.Marshal(messageData)
		if err != nil {
			apiErr := errors.InternalServer("序列化请求数据失败")
			errors.WriteJSON(w, apiErr)
			return
		}

		if err := json.Unmarshal(jsonData, &msg); err != nil {
			apiErr := errors.BadRequest("无效的广播消息格式")
			errors.WriteJSON(w, apiErr)
			return
		}

		// 验证 Channel 是否有效
		if msg.Channel != nil {
			channel := *msg.Channel
			if channel != News && channel != Alert && channel != Promotion {
				apiErr := errors.BadRequest("无效的广播渠道")
				errors.WriteJSON(w, apiErr)
				return
			}
		}

		// 这里可以添加创建广播消息的业务逻辑
		msg.Id = getNextID() // 模拟ID生成
		msg.Status = Sending // 设置初始状态
		responseMsg = msg

	default:
		apiErr := errors.BadRequest(fmt.Sprintf("不支持的消息类型: %s", messageType))
		errors.WriteJSON(w, apiErr)
		return
	}

	// 返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(responseMsg)
}

// FindMessagesByNumber 根据手机号查询短信消息
func FindMessagesByNumber(w http.ResponseWriter, r *http.Request) {
	numberParam := chi.URLParam(r, "number")

	// 验证手机号
	_, err := strconv.ParseInt(numberParam, 10, 64)
	if err != nil {
		apiErr := errors.BadRequest("无效的手机号格式")
		errors.WriteJSON(w, apiErr)
		return
	}

	// 这里应该添加根据手机号查询数据的逻辑
	// 模拟数据
	messages := []Message{
		{
			Id:      1,
			Content: "验证码：929253，有效期10分钟。如非本人操作，请忽略。",
			Type:    Sms,
			Status:  Sent,
		},
		{
			Id:      2,
			Content: "您的订单已发货，物流单号：SF12345678",
			Type:    Sms,
			Status:  Sent,
		},
	}

	response := MessageListResponse{
		Messages: &messages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// FindMessagesByUID 根据用户ID查询站内消息
func FindMessagesByUID(w http.ResponseWriter, r *http.Request) {
	uidParam := chi.URLParam(r, "uid")

	// 验证用户ID
	uid, err := strconv.ParseInt(uidParam, 10, 64)
	if err != nil {
		apiErr := errors.BadRequest("无效的用户ID格式")
		errors.WriteJSON(w, apiErr)
		return
	}

	// 这里应该添加根据用户ID查询数据的逻辑
	// 模拟数据
	messages := []Message{
		{
			Id:      3,
			Content: "您的账户已成功升级为VIP会员",
			Type:    Sitemessage,
			Status:  Read,
		},
		{
			Id:      4,
			Content: "系统将于明日凌晨2点进行维护，预计持续2小时",
			Type:    Sitemessage,
			Status:  Received,
		},
	}

	// 打印用户ID，确保被正确解析（实际生产环境可删除）
	fmt.Printf("查询用户ID: %d\n", uid)

	response := MessageListResponse{
		Messages: &messages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// 生成下一个消息ID的辅助函数（实际应用中应该由数据库生成）
var lastID int64 = 100

func getNextID() int64 {
	lastID++
	return lastID
}

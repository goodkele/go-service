package userbase

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/giant-tech/go-service/base/imsg"
	"github.com/giant-tech/go-service/framework/entity"
	"github.com/giant-tech/go-service/framework/idata"

	assert "github.com/aurelien-rainone/assertgo"
	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
	"go.uber.org/atomic"
)

// GateUserBase 基础网关信息
type GateUserBase struct {
	entity.Entity
	isClosed  atomic.Bool //玩家协程是否已经关闭
	isLogout  bool        //是否已经调用Logout
	isRunning bool
}

// OnEntityLoop lobby user的逻辑帧
func (gu *GateUserBase) OnEntityLoop() {
	gu.FlushDirtyProp()

	gu.Entity.OnEntityLoop()
}

// Run 运行玩家协程
func (gu *GateUserBase) Run() {
	assert.False(gu.isRunning, "gate user Run() again")
	gu.isRunning = true

	ticker := time.NewTicker(time.Duration(50) * time.Millisecond)
	defer ticker.Stop()

	//属性保存到数据库间隔时间
	propsSaveTicker := time.NewTicker(time.Duration(10) * time.Second)
	defer propsSaveTicker.Stop()

	// 调试状态, 向客户端发送服务器错误信息
	isDebug := viper.GetBool("Config.Debug")

	doActionFunc := func() {
		defer func() {
			if err := recover(); err != nil {
				strErr := fmt.Sprintf("%v", err)
				strStack := string(debug.Stack())
				log.Error("GateUserBase: ", strErr, strStack)
				if isDebug {
					gu.AsyncCall(idata.ServiceGateway, "ServerError", strErr, strStack)
				}
			}
		}()

		select {
		case <-ticker.C:
			gu.MainLoop()

		case data := <-gu.DataC:
			gu.ProcessCall(data)
		}
	}

	for !gu.isClosed.Load() {
		doActionFunc()
	}

	log.Debug("gate user goroutine end, uid: ", gu.GetEntityID())
}

//CloseRoutine 关闭协程
func (gu *GateUserBase) CloseRoutine() {
	//延迟几秒关闭玩家协程
	time.AfterFunc(5*time.Second, func() {
		gu.isClosed.Store(true)
	})
}

//CloseCliSession 关闭客户端连接
func (gu *GateUserBase) CloseCliSession() {
	if gu.CliSess != nil {
		gu.CliSess.Close()
		gu.CliSess = nil
	}
}

//Logout 登出
func (gu *GateUserBase) Logout() {
	if gu.isLogout {
		return
	}

	gu.isLogout = true

	log.Info("Logout(), id: ", gu.GetEntityID())

	//loginclt.GetLoginCliMgr().PlayerLogout(gu.GetEntityID())

	// 关闭客户端连接
	gu.CloseCliSession()

	gu.GetIEntities().DestroyEntity(gu.GetEntityID())
}

// AsyncSend 由LobbyUserBase协程发送
func (gu *GateUserBase) AsyncSend(msg imsg.IMsg) {
	gu.PostFunction(func() {
		gu.Send(msg)
	})
}

// AsyncSendRaw 由LobbyUserBase协程发送
func (gu *GateUserBase) AsyncSendRaw(buff []byte) {
	gu.PostFunction(func() {
		gu.SendRaw(buff)
	})
}

// Send 发送消息
func (gu *GateUserBase) Send(msg imsg.IMsg) error {
	if gu.CliSess != nil {
		return gu.CliSess.Send(msg)
	}

	return fmt.Errorf("CliSess is nil")
}

// SendRaw 发送原始数据
func (gu *GateUserBase) SendRaw(buff []byte) error {
	if gu.CliSess != nil {
		return gu.CliSess.SendRaw(buff)
	}

	return fmt.Errorf("CliSess is nil")
}

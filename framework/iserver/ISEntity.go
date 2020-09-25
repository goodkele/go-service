package iserver

import (
	"github.com/giant-tech/go-service/base/imsg"
	"github.com/giant-tech/go-service/base/linmath"
)

// ICellEntity 代表空间中的一个实体
type ICellEntity interface {
	IEntity
	IPos
	IClientBroadcaster
	IAOIEntity
	//SetClient(ISess)
}

// ICellEntityPropsGetter 获取属性数据
type ICellEntityPropsGetter interface {
	GetAOIProp() (int, []byte)
}

// IPos 拥有位置信息的接口
type IPos interface {
	SetPos(pos linmath.Vector3)
	GetPos() linmath.Vector3

	SetRota(linmath.Vector3)
	GetRota() linmath.Vector3
}

// IPosChange 位置变化回调
type IPosChange interface {
	OnPosChange(curPos, origPos linmath.Vector3)
}

// IMover 寻路移动能力
type IMover interface {
	SetSpeed(speed float32)
	Move(destPos linmath.Vector3) bool
	StopMove()
	IsMoving() bool

	PauseNav()
	ResumeNav()
}

// IClientBroadcaster AOI广播
type IClientBroadcaster interface {
	CastMsgToAllClient(imsg.IMsg)
	CastMsgToMe(imsg.IMsg)
	CastMsgToAllClientExceptMe(imsg.IMsg)
	CastMsgToRangeExceptMe(center *linmath.Vector3, radius int, msg imsg.IMsg)
	CastMsgToCenterExceptMe(center *linmath.Vector3, radius int, msg imsg.IMsg)

	CastRPCToAllClient(methodName string, args ...interface{})
	CastRPCToMe(methodName string, args ...interface{})
	CastRPCToAllClientExceptMe(methodName string, args ...interface{})
}

// IAOIEntity  AOI实体类型查询
type IAOIEntity interface {
	IsWatcher() bool
	// IsMarker() bool
}

// ISyncToGhosts 同步ghost接口
type ISyncToGhosts interface {
	SyncToGhosts(imsg.IMsg)
}

// ISendMsgToReal 发消息到real接口
type ISendMsgToReal interface {
	SendMsgToReal(msg imsg.IMsg)
}

// AOIInfo aoi信息
type AOIInfo struct {
	IsEnter bool
	Entity  ICoordEntity
}

// IAOIUpdate  AOI更新，由逻辑层实现
type IAOIUpdate interface {
	AOIUpdate([]AOIInfo)
}

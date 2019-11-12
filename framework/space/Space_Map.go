package space

import (
	"errors"
	"fmt"
	"math"
	"zeus/linmath"

	"github.com/cihub/seelog"
)

// SetMap 设置地图
func (s *Space) SetMap(mapName string) {
	s.mapName = mapName
}

func (s *Space) loadMap() {
	ret := s._loadMap()
	s.onLoadMapFinished(ret)
}

func (s *Space) _loadMap() bool {
	task := getMapsInst().getMap(s.mapName)
	ret := <-task.c

	//TODO fixme 临时处理
	//if ret {
	s.mapInfo = task.info
	//}

	task.c <- ret
	return ret
}

func (s *Space) onLoadMapFinished(ret bool) {
	im, ok := s.GetRealPtr().(IMapLoader)

	seelog.Debugf("space %p onMap finished %v ", s, ret)

	if ok {
		if ret {
			im.OnMapLoadSucceed()
		} else {
			im.OnMapLoadFailed()
		}
	}

	s.isMapLoaded = ret
}

// IsMapLoaded 地图是否加载
func (s *Space) IsMapLoaded() bool {
	return s.isMapLoaded
}

// FindPath 查询一条路径
func (s *Space) FindPath(srcPos, destPos linmath.Vector3) ([]linmath.Vector3, error) {

	if s.mapInfo == nil {
		return nil, errors.New("no map info")
	}

	return s.mapInfo.pathFinder.FindPath(srcPos, destPos)
}

// Raycast 以origin为原点, direction为方向, length为长度, 作射线检测, mask为射线检测的层级
func (s *Space) Raycast(origin, direction linmath.Vector3, length float32, mask int32) (float32, linmath.Vector3, int32, bool, error) {
	if s.mapInfo == nil {
		return 0, linmath.Vector3_Invalid(), 0, false, fmt.Errorf("no map info")
	}

	dist, pos, layer, hit := s.mapInfo.scene.Raycast(origin, direction, length, mask)
	return dist, pos, layer, hit, nil
}

// CapsuleRaycast 胶囊体检测
func (s *Space) CapsuleRaycast(head, foot linmath.Vector3, radius float32, origin, direction linmath.Vector3, length float32) (float32, bool, error) {
	if s.mapInfo == nil {
		return math.MaxFloat32, false, fmt.Errorf("no map info when CapsuleRaycast")
	}

	dist, hit := s.mapInfo.scene.CapsuleRaycast(head, foot, radius, origin, direction, length)
	return dist, hit, nil
}

func (s *Space) SphereRaycast(center linmath.Vector3, radius float32, origin, direction linmath.Vector3, length float32) (float32, bool, error) {
	if s.mapInfo == nil {
		return math.MaxFloat32, false, fmt.Errorf("no map info when SphereRaycast")
	}

	dist, hit := s.mapInfo.scene.SphereRaycast(center, radius, origin, direction, length)
	return dist, hit, nil
}

// GetRanges 获取区域管理器
func (s *Space) GetRanges() *MapRanges {
	return s.mapInfo.ranges
}

// GetHeight 获取高度
func (s *Space) GetHeight(x, z float32) (float32, error) {
	origin := linmath.Vector3{
		X: x,
		Y: 1000,
		Z: z,
	}
	direction := linmath.Vector3{
		X: 0,
		Y: -1,
		Z: 0,
	}
	_, pos, _, hit := s.mapInfo.scene.Raycast(origin, direction, 2000, 1<<12)
	if hit {
		return pos.Y, nil
	}
	return 0, fmt.Errorf("射线检测失败, 无法获取高度")
}

// IsWater 判断坐标点是否是水域
func (s *Space) IsWater(x, z float32, waterlevel float32) (bool, error) {
	height, err := s.GetHeight(x, z)
	if err != nil {
		return false, err
	}

	return height <= waterlevel, nil
}

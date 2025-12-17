package cron

import "time"

// ConstantDelaySchedule 表示一个简单的重复执行周期，例如"每5分钟一次"
// 它不支持频率高于每秒一次的作业
type ConstantDelaySchedule struct {
	Delay time.Duration
}

// Every 返回一个按指定持续时间激活的 crontab 计划
//   - 不支持小于一秒的延迟（将向上舍入到1秒）
//   - 任何小于一秒的字段都将被截断
//   - duration: 持续时间
//
// 返回值: ConstantDelaySchedule 结构体
func Every(duration time.Duration) ConstantDelaySchedule {
	if duration < time.Second {
		duration = time.Second
	}
	return ConstantDelaySchedule{
		Delay: duration - time.Duration(duration.Nanoseconds())%time.Second,
	}
}

// Next 返回下一次应该运行的时间
//   - 这会进行四舍五入，使得下次激活时间将在整秒处
//   - t: 参考时间
//
// 返回值: 下次运行的时间
func (schedule ConstantDelaySchedule) Next(t time.Time) time.Time {
	return t.Add(schedule.Delay - time.Duration(t.Nanosecond())*time.Nanosecond)
}

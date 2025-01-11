package helper

//var BEGIN_TIMESTAMP int64 = 1672531200
//
//// 序列号的bit位
//var COUNT_BITS int64 = 32
//
//func NextId(keyPrefix string) int64 {
//	//1.时间戳
//	now_timestamp := int64(time.Now().Unix())
//	diff_timestamp := now_timestamp - BEGIN_TIMESTAMP
//	//2.自增长序列号
//	now_date := time.Now().Format("20060102")
//	count, _ := models.RDb.Incr(context.Background(), define.ICR_KEY+keyPrefix+":"+now_date).Result()
//
//	return diff_timestamp<<COUNT_BITS | count
//}

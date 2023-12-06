package redis_test

// var (
// 	RedisHost = "localhost"
// 	RedisPort = "6379"
// 	options   redis.Options
// )

// func init() {
// 	host, exists := os.LookupEnv("REDIS_HOST")
// 	if exists {
// 		RedisHost = host
// 	}
// 	port, exists := os.LookupEnv("REDIS_PORT")
// 	if exists {
// 		RedisPort = port
// 	}
// 	options = redis.Options{
// 		Address:  RedisHost + ":" + RedisPort,
// 		Password: "",
// 		DB:       15,
// 	}
// }

// func TestClose(t *testing.T) {
// 	s, _ := test.MustCreateStorage[int](t, "redis", options)

// 	test.Close(t, s)
// }

// func TestSetGetDelete(t *testing.T) {
// 	s, closer := test.MustCreateStorage[int](t, "redis", options)
// 	defer closer()

// 	test.SetGetDelete(t, s)
// }

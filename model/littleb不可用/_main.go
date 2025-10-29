package main

/*
func main() {

		out, err := exec.Command("ping", "www.baidu.com").Output()
		if err != nil {
			log.Println(err)
		} else {
			log.Println(string(out))
		}
	}

	func fibonacci(c, quit chan int) {
		x, y := 1, 1
		for {
			select {
			case c <- x:
				x, y = y, x+y
			case <-quit:
				fmt.Println("quit")
				return
			}
		}
	}
*/

func main() {

	/*
		c := make(chan int)
		quit := make(chan int)
		go func() {
			for i := 0; i < 10; i++ {
				fmt.Println(<-c)
			}
			quit <- 0
		}()
		fibonacci(c, quit)
	*/
}

package main

type Letter struct {
	Text string `form:"text" binding:"required"`
	Ttl  int    `form:"ttl" binding:"required"`
}

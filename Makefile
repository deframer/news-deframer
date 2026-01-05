ifneq ("$(wildcard .env)","")
  #$(info using .env file)
  include .env
  export $(shell sed 's/=.*//' .env)
endif

all:
	@echo all


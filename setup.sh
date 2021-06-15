#!/bin/bash

if [ ! -f .env ]; then
    cat .env.example > .env
fi
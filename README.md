# GO GIF Generator API

## Overview
The GO API allows a client to generate a GIF file out of a video file.

The conversion process is based on the [gocv](https://gocv.io/) library (GO's variant for [opencv](https://opencv.org/))

In order to maintain efficiency the API is using GOD's built-in  [workers pool](https://gobyexample.com/worker-pools).


## Requests

### /gifgen
Trigger the conversion process. Provides the server with a video file.
The request returns a gif file if the process has been completed successfully. 

### /test
An endpoint for helath check
# Multiple streams example

This example represents how [source](https://docs.yomo.run/source) pipe the local file stream to [zipper](https://docs.yomo.run/zipper), and zipper pipe the stream to multiple [stream functions](https://docs.yomo.run/stream-fn).

## Code structure

+ `source`: Read the files in a local directory, create a new QUIC stream for each file to pipe the file stream. [docs.yomo.run/source](https://docs.yomo.run/source)
+ `zipper`: Receive the data from `source`, create a new stream and pipe the `source` stream to `stream-fn` [docs.yomo.run/zipper](https://docs.yomo.run/zipper)
+ `sfn`: Receive the stream from `zipper` and store the file in local via `io.copy`. [docs.yomo.run/stream-function](https://docs.yomo.run/stream-fn)

## Prepare

Install YoMo CLI

### Binary (Recommended)

```bash
$ curl -fsSL "https://bina.egoist.sh/yomorun/cli?name=yomo" | sh

  ==> Resolved version latest to v1.2.1
  ==> Downloading asset for darwin amd64
  ==> Installing yomo to /usr/local/bin
  ==> Installation complete
```

## Run

### Run [zipper](https://docs.yomo.run/zipper)

```bash
cd ./zipper
go run main.go

2022-09-07 15:28:18.324	[yomo:zipper] Listening SIGUSR1, SIGUSR2, SIGTERM/SIGINT...
2022-09-07 15:28:18.328	[core:server] ✅ [Zipper][17916] Listening on: 127.0.0.1:9000, MODE: DEVELOPMENT, QUIC: [v1 draft-29], AUTH: [none]
```

### Run [sfn-1](https://docs.yomo.run/stream-fn)

```bash
cd ./sfn
go run main.go sink-1

2022-09-07 15:29:17.626	[core:client] use credential: [none]
2022-09-07 15:29:17.632	[core:client] ❤️  [sink-1][jabXFp5WpHDin5o-mYaId]([::]:61242) is connected to YoMo-Zipper localhost:9000
```

### Run [sfn-2](https://docs.yomo.run/stream-fn)
```bash
cd ./sfn
go run main.go sink-2

2022-09-07 15:29:51.884	[core:client] use credential: [none]
2022-09-07 15:29:51.890	[core:client] ❤️  [sink-2][aJb-JKHYYytNq_V2JEJbF]([::]:62413) is connected to YoMo-Zipper localhost:9000
```

### Run [yomo-source](https://docs.yomo.run/source)

```bash
cd ./source
go run main.go /path/to/dir

2022-09-07 15:30:46.810	[core:client] use credential: [none]
2022-09-07 15:30:46.815	[core:client] ❤️  [source][nrxQzDFtSAr6a5oPJRCSk]([::]:58333) is connected to YoMo-Zipper localhost:9000
```

### Results

The terminal of `yomo-srouce` will print the real-time receives value.

```bash

2022-09-07 15:30:47.817	sending file test1.mp4 to yomo-zipper...
2022-09-07 15:30:47.817	sending file test2.mp4 to yomo-zipper...
2022-09-07 15:31:07.196	file: test1.mp4, written: 366676386
2022-09-07 15:31:07.196	file: test1.mp4, md5: a6a87007a45e7d35846adb11c118ee1d
2022-09-07 15:31:09.027	file: test2.mp4, written: 434894207
2022-09-07 15:31:09.027	file: test2.mp4, md5: 372148e7d1ba577914047a1ec4580dc9
```

The terminal of `sfn-1` will print the real-time noise value.

```bash
2022-09-07 15:30:47.819	receiving file: test1.mp4
2022-09-07 15:30:47.819	receiving file: test2.mp4
2022-09-07 15:31:07.218	written: 366676386, /path/to/dir/sink-1-test1.mp4
2022-09-07 15:31:07.218	file: test1.mp4, md5: a6a87007a45e7d35846adb11c118ee1d
2022-09-07 15:31:09.062	written: 434894207, /path/to/dir/sink-1-test2.mp4
2022-09-07 15:31:09.062	file: test2.mp4, md5: 372148e7d1ba577914047a1ec4580dc9
```

The terminal of `sfn-2` will print the real-time noise value.

```bash
2022-09-07 15:30:47.820	receiving file: test1.mp4
2022-09-07 15:30:47.820	receiving file: test2.mp4
2022-09-07 15:31:07.218	written: 366676386, /path/to/dir/sink-2-test1.mp4
2022-09-07 15:31:07.218	file: test1.mp4, md5: a6a87007a45e7d35846adb11c118ee1d
2022-09-07 15:31:09.062	written: 434894207, /path/to/dir/sink-2-test2.mp4
2022-09-07 15:31:09.062	file: test2.mp4, md5: 372148e7d1ba577914047a1ec4580dc9
```

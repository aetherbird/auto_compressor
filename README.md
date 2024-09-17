# auto_compressor.go

A command-line tool for Linux, designed to compress MP4 video files to a desired output size using FFmpeg. Calculates the appropriate video bitrate based on the duration of the video and compresses it to the target size while maintaining audio quality.

## Dependencies

go and ffmpeg

## Installation

If you are on a Debian, Ubuntu, or Mint-based distro, you can use the following one-liner to install dependencies and clone the script:

```bash
curl -s https://raw.githubusercontent.com/aetherbird/auto_compressor/main/auto_compressor_installer.sh | bash
```

## How to Use

Once you've installed go and cloned the auto_compressor program, you can use the following command to compress your videos:

```bash
go run auto_compressor.go <input_file.mp4> <desired_output_size_MB>
```

For example, to compress a video to 50 MB:

```bash
go run auto_compressor.go ~/Videos/my_video.mp4 50
```


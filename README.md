# Wallpaper Generator Tool

This tool provides a command line interface to create many different snapshots of images.
It is primarily a wrapper around ImageMagick that carries out specific set of actions against provided images.

## Help

```
Manipulate images for use as desktop wallpapers

Usage:
  wp [flags]
  wp [command]

Available Commands:
  extract     Extract image slices
  help        Help about any command
  pick        Pick a single image slice

Flags:
  -h, --help      help for wp
  -v, --version   Print the application version and exit

Use "wp [command] --help" for more information about a command.
```

## Usage

### Extract

Extracts a bunch of sub-images from a given source image, and writes them to the provided directory.
Buckets images by resolution.

```
$ mkdir images
$ wp extract 1024x768 images https://i.imgur.com/hqCBTK8.png
/path/to/images/1024x768/hqCBTK8_scaled_west.png
/path/to/images/1024x768/hqCBTK8_scaled_center.png
/path/to/images/1024x768/hqCBTK8_scaled_east.png
/path/to/images/1024x768/hqCBTK8_north.png
/path/to/images/1024x768/hqCBTK8_northeast.png
/path/to/images/1024x768/hqCBTK8_east.png
/path/to/images/1024x768/hqCBTK8_southeast.png
/path/to/images/1024x768/hqCBTK8_south.png
/path/to/images/1024x768/hqCBTK8_southwest.png
/path/to/images/1024x768/hqCBTK8_west.png
/path/to/images/1024x768/hqCBTK8_northwest.png
/path/to/images/1024x768/hqCBTK8_center.png
$ tree images
images
└── 1024x768
    ├── hqCBTK8_center.png
    ├── hqCBTK8_east.png
    ├── hqCBTK8_north.png
    ├── hqCBTK8_northeast.png
    ├── hqCBTK8_northwest.png
    ├── hqCBTK8_scaled_center.png
    ├── hqCBTK8_scaled_east.png
    ├── hqCBTK8_scaled_west.png
    ├── hqCBTK8_south.png
    ├── hqCBTK8_southeast.png
    ├── hqCBTK8_southwest.png
    └── hqCBTK8_west.png

1 directory, 12 files
```

Each of these images can be evaluated for being optimal for your use case.
Once you've chosen your favorite, use `pick` moving forward.

### Pick

Selects a single predetermined image to extract from a given source image.
Buckets images by resolution.

```
$ mkdir images
$ wp pick 1024x768 images west --scaled https://i.imgur.com/hqCBTK8.png
/path/to/images/1024x768/hqCBTK8_scaled_west.png
$ tree images
images
└── 1024x768
    └── hqCBTK8_scaled_west.png

1 directory, 1 file
```

This operation could be used to fill up an entire directory of preferred wallpapers.

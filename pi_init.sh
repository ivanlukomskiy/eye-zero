sudo apt update && sudo apt full-upgrade -y

# download and unpack the lib
wget https://gstreamer.freedesktop.org/src/gstreamer/gstreamer-1.18.4.tar.xz
sudo tar -xf gstreamer-1.18.4.tar.xz
cd gstreamer-1.18.4
# make an installation folder
mkdir build && cd build
# run meson (a kind of cmake)
meson --prefix=/usr \
        --wrap-mode=nofallback \
        -D buildtype=release \
        -D gst_debug=true \
        -D package-origin=https://gstreamer.freedesktop.org/src/gstreamer/ \
        -D package-name="GStreamer 1.18.4 BLFS" ..

# build the software
ninja -j4
# test the software (optional)
#ninja test
# install the libraries
sudo ninja install
sudo ldconfig


cd ~
# download and unpack the plug-ins base
wget https://gstreamer.freedesktop.org/src/gst-plugins-base/gst-plugins-base-1.18.4.tar.xz
sudo tar -xf gst-plugins-base-1.18.4.tar.xz
# make an installation folder
cd gst-plugins-base-1.18.4
mkdir build
cd build
# run meson
meson --prefix=/usr \
-D buildtype=release \
-D package-origin=https://gstreamer.freedesktop.org/src/gstreamer/ ..
ninja -j4
# optional
#ninja test
# install the libraries
sudo ninja install
sudo ldconfig


cd ~
sudo apt-get install libjpeg-dev
# download and unpack the plug-ins good
wget https://gstreamer.freedesktop.org/src/gst-plugins-good/gst-plugins-good-1.18.4.tar.xz
sudo tar -xf gst-plugins-good-1.18.4.tar.xz
cd gst-plugins-good-1.18.4
# make an installation folder
mkdir build && cd build
# run meson
meson --prefix=/usr \
       -D buildtype=release \
       -D package-origin=https://gstreamer.freedesktop.org/src/gstreamer/ \
       -D package-name="GStreamer 1.18.4 BLFS" ..
ninja -j4
# optional
#ninja test
# install the libraries
sudo ninja install
sudo ldconfig


cd ~
sudo apt install -y librtmp-dev libvo-aacenc-dev
wget https://gstreamer.freedesktop.org/src/gst-plugins-bad/gst-plugins-bad-1.18.4.tar.xz
sudo tar -xf gst-plugins-bad-1.18.4.tar.xz
cd gst-plugins-bad-1.18.4
mkdir build && cd build
meson --prefix=/usr \
       -D buildtype=release \
       -D package-origin=https://gstreamer.freedesktop.org/src/gstreamer/ \
       -D package-name="GStreamer 1.18.4 BLFS" ..
ninja -j4
sudo ninja install
sudo ldconfig


cd ~
wget https://gstreamer.freedesktop.org/src/gst-plugins-ugly/gst-plugins-ugly-1.18.4.tar.xz
sudo tar -xf gst-plugins-ugly-1.18.4.tar.xz
cd gst-plugins-ugly-1.18.4
mkdir build && cd build
meson --prefix=/usr \
      -D buildtype=release \
      -D package-origin=https://gstreamer.freedesktop.org/src/gstreamer/ \
      -D package-name="GStreamer 1.18.4 BLFS" ..
ninja -j4
sudo ninja install
sudo ldconfig


cd ~
wget https://gstreamer.freedesktop.org/src/gst-omx/gst-omx-1.18.4.tar.xz
sudo tar -xf gst-omx-1.18.4.tar.xz
cd gst-omx-1.18.4
mkdir build && cd build
meson --prefix=/usr \
       -D header_path=/opt/vc/include/IL \
       -D target=rpi \
       -D buildtype=release ..
ninja -j4
sudo ninja install
sudo ldconfig


cd ~
wget https://gstreamer.freedesktop.org/src/gst-rtsp-server/gst-rtsp-server-1.18.4.tar.xz
tar -xf gst-rtsp-server-1.18.4.tar.xz
cd gst-rtsp-server-1.18.4
mkdir build && cd build
meson --prefix=/usr
       --wrap-mode=nofallback
       -D buildtype=release
       -D package-origin=https://gstreamer.freedesktop.org/src/gstreamer/
       -D package-name="GStreamer 1.18.4 BLFS" ..
ninja -j4
sudo ninja install
sudo ldconfig
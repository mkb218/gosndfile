gosndfile is a binding for [libsndfile][1]. It is distributed under the same terms (your choice of LGPL 2.1 or 3). If you install libsndfile outside of your system include and lib paths, make sure to set the environment variable PKG_CONFIG_PATH accordingly. This package should be go get-able: e.g. `go get github.com/mkb218/gosndfile/sndfile`

Tests will only pass on libsndfile 1.0.26.

   [1]: http://www.mega-nerd.com/libsndfile/

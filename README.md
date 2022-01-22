

### Installation
> Tested on macos M1


```
brew install openssl
```

```
git clone git@github.com:bitcoin-core/secp256k1.git
cd secp256k1
./autogen.sh
./configure --enable-module-recovery

make
sudo make install

```

If on macos M1
```
export PATH="/opt/homebrew/opt/openssl@3/bin:$PATH"
export CPPFLAGS="-I/opt/homebrew/opt/openssl@3/include"
export LDFLAGS="-L/opt/homebrew/opt/openssl@3/lib"


git clone git@github.com:dfoxfranke/libaes_siv.git
cd libaes_siv
cmake -DCMAKE_PREFIX_PATH=/opt/homebrew/opt/openssl@3 .
make
sudo make install

```



```
export PKG_CONFIG_PATH="/opt/homebrew/opt/openssl@3/lib/pkgconfig"
cd urcrypt/
./autogen.sh
./configure --disable-shared
make
sudo make install
```

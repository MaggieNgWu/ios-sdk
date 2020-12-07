# go/Makefile

##ios-arm64:
 CGO_ENABLED=1 \
 GOOS=darwin \
 GOARCH=arm64 \
 SDK=iphoneos \
 CC=$(PWD)/clangwrap.sh \
 CGO_CFLAGS="-fembed-bitcode" \
 go build -buildmode=c-archive -tags ios -o ./build/arm64.a .

##ios-x86_64:
 CGO_ENABLED=1 \
 GOOS=darwin \
 GOARCH=amd64 \
 SDK=iphonesimulator \
 CC=$(PWD)/clangwrap.sh \
 go build -buildmode=c-archive -tags ios -o ./build/x86_64.a .

##ios: ios-arm64 ios-x86_64
lipo ./build/x86_64.a ./build/arm64.a -create -output ./build/go-sdk-ios.a
cp ./build/arm64.h ./build/go-sdk-ios.h
cp data_struct.h ./build/data_struct.h
rm ./build/arm64.h ./build/arm64.a ./build/x86_64.a ./build/x86_64.h

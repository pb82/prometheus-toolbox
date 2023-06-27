#!/usr/bin/env sh

archs=(amd64 arm64)
oses=(linux darwin)

for arch in "${archs[@]}"
do
  for os in "${oses[@]}"
  do
    env GOOS="${os}" GOARCH="${arch}" go build -o prometheus-toolbox_"${os}"_"${arch}"
  done
done


# image-processing-app
Склонируй репозиторий
git clone https://github.com/IslomSobirov/image-processing-app


## Собери Docker-образ
``` docker build -t image-processor . ```


## Запусти контейнер
``` 
docker run -d -p 8085:8085 \
  -v $(pwd)/img_orig:/img_orig \
  -v $(pwd)/img_res:/img_res \
  image-processor \
  ./app -path-orig /img_orig -path-res /img_res -width 300 -height 300
```

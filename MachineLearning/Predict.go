package MachineLearning

import (
	"GoFender/Utils"
	"GoFender/YamlConfig"
	"fmt"
	"github.com/corona10/goimagehash"
	tf "github.com/galeone/tensorflow/tensorflow/go"
	tg "github.com/galeone/tfgo"
	tgi "github.com/galeone/tfgo/image"
	"image"
	"image/png"
	"log"
	"os"
)

func max(res []float32) (int, float32) {
	maxIndex := 0
	maxVal := res[0]
	for i := 0; i < len(res); i++ {
		if maxVal < res[i] {
			maxVal = res[i]
			maxIndex = i
		}
	}
	return maxIndex, maxVal
}

func encodePacket(packet Utils.CommonPacket) (string, error) {
	ImgData := make([]byte, 28*28)
	data := packet.ComPacketData
	copy(ImgData, data)
	if len(data) <= 28*28 {
		for i := len(data); i <= 28*28-1; i++ {
			ImgData[i] = 0x00
		}
	}
	img := image.NewGray(image.Rect(0, 0, 28, 28))
	img.Pix = ImgData
	hash, err := goimagehash.AverageHash(img)
	if err != nil {
		log.Fatal("Error computing hash:", err)
	}
	file, err := os.Create(YamlConfig.Myconfig.TempPath + fmt.Sprintf("%x", hash.ToString()) + ".png")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	err = png.Encode(file, img)
	if err != nil {
		log.Fatal(err)
	}
	return file.Name(), nil
}

func predict(imgname string) (int, float32) {
	model := tg.LoadModel(YamlConfig.Myconfig.ModelPath, []string{"serve"}, nil)
	root := tg.NewRoot()
	img := tgi.Read(root, imgname, 1)
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			log.Fatal(err)
		}
	}(imgname)
	restf := tg.Exec(root, []tf.Output{img.Value()}, nil, &tf.SessionOptions{})
	xInput, _ := tf.NewTensor(restf[0].Value())
	err := xInput.Reshape([]int64{1, 28, 28, 1})
	if err != nil {
		panic(err)
	}
	results := model.Exec([]tf.Output{
		model.Op("StatefulPartitionedCall", 0),
	}, map[tf.Output]*tf.Tensor{
		model.Op("serving_default_inputs_input", 0): xInput,
	})
	res := results[0].Value().([][]float32)[0]
	return max(res)
}

func Attacktype(packet Utils.CommonPacket) (float64, string) {
	var acType []string
	acType = append(acType, "Cridex", "Geodo", "Htbot", "Miuref", "Neris", "Shifu", "Tinba", "Virut", "Zeus")
	imgname, err := encodePacket(packet)
	if err != nil {
		panic(err)
	}
	in, value := predict(imgname)
	return float64((value * 100) * 0.80), acType[in]
}

## 模型存储位置

该文件夹存储模型文件，模型使用TensorFlow SavedModel格式存储。

模型文件夹结构如下：

```
model
├── keras
│   ├── assets
│   ├── saved_model.pb
│   └── variables
│       ├── variables.data-00000-of-00001
│       └── variables.index
```

该模型训练使用Python，TensorFlow
2.0版本，模型设计逻辑及训练代码参考[BeStrongok/Malicious-Traffic-Classification](https://github.com/BeStrongok/Malicious-Traffic-Classification)

模型训练好之后需要导出为PB格式，之后放到该文件夹下，即可使用。


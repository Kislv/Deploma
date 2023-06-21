import os
import io
import sys
from PIL import Image

sys.path.append( '../ML/DiseaseClassification' )
import implementModel

IMAGES = 'Images'

IMAGE_CLASSIFICATION_PATH = os.path.join(IMAGES, 'mole_pores2.jpg')
readed_image = Image.open(IMAGE_CLASSIFICATION_PATH)
imgByteArr = io.BytesIO()
readed_image.save(imgByteArr, format=readed_image.format)
imgByteArr = imgByteArr.getvalue()
prediction = implementModel.classificateImage(".jpg", (1,1), imgByteArr)
assert "Vasculitis Photos" == prediction,  "Wrong answer!"


IMAGE_CLASSIFICATION_PATH = os.path.join(IMAGES, 'mole_pores1.jpg')
readed_image = Image.open(IMAGE_CLASSIFICATION_PATH)
imgByteArr = io.BytesIO()
readed_image.save(imgByteArr, format=readed_image.format)
imgByteArr = imgByteArr.getvalue()
prediction = implementModel.classificateImage(".jpg", (1,1), imgByteArr)
assert "Nail Fungus and other Nail Disease" == prediction,  "Wrong answer!"

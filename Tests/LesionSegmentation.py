import os
import io
import sys
from PIL import Image

sys.path.append( '../ML/LesionSegmentation' )
import implementModel

IMAGES = 'Images'

IMAGE_CLASSIFICATION_PATH = os.path.join(IMAGES, 'mole_pores2.jpg')
readed_image = Image.open(IMAGE_CLASSIFICATION_PATH)
imgByteArr = io.BytesIO()
readed_image.save(imgByteArr, format=readed_image.format)
imgByteArr = imgByteArr.getvalue()
segmented = implementModel.segment_lesion(".jpg", (1,1), imgByteArr)
assert type(segmented) == type(imgByteArr) and segmented != imgByteArr,  "Wrong answer!"


IMAGE_CLASSIFICATION_PATH = os.path.join(IMAGES, 'mole_pores1.jpg')
readed_image = Image.open(IMAGE_CLASSIFICATION_PATH)
imgByteArr = io.BytesIO()
readed_image.save(imgByteArr, format=readed_image.format)
imgByteArr = imgByteArr.getvalue()
segmented = implementModel.segment_lesion(".jpg", (1,1), imgByteArr)
assert type(segmented) == type(imgByteArr) and segmented != imgByteArr,  "Wrong answer!"

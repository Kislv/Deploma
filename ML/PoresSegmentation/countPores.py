import math
import cv2
from skimage.transform import resize
from skimage.morphology import binary_erosion
from matplotlib import pyplot as plt
from skimage.io import imread
from skimage.measure import label, regionprops
from skimage.color import label2rgb
from PIL import Image
import grpc
import sys
import cv2 as cv

import os
import io
import numpy as np

sys.path.append( '../../GRPC/LesionSegmentation/ML' )
import lesionSegmentation_pb2
import lesionSegmentation_pb2_grpc

ROOT = '/home/viktor/Projects/Univer/Sem8/Deploma/Deploma'
FOLDER_PATH = 'ML/PoresSegmentation'
IMAGES_PATH = 'Images'
BASE_IMAGE_PATH = 'mole_pores2.jpg'
LESION_BW_PATH = 'after.png'

def save_image(image, file_name, is_gray = False ):
  """
  Save numpy array as png.

  :param image: numpy array
  :param file_name: name of file 
  :param is_gray: is this image gray
  :return: returns nothing
  """
  cmap = None
  if is_gray:
    cmap = 'gray'
  plt.imshow(image, interpolation='nearest', cmap = cmap)
  plt.axis('off')
  plt.savefig(os.path.join(IMAGES_PATH, file_name), bbox_inches='tight', pad_inches = 0)

def CLAHE_to_image (image):
  bgr = image

  lab = cv2.cvtColor(bgr, cv2.COLOR_BGR2LAB)

  lab_planes = cv2.split(lab)
  clahe = cv2.createCLAHE(clipLimit=2.0,tileGridSize=(8,8))

  lab_planes = list(lab_planes)

  lab_planes[0] = clahe.apply(lab_planes[0])

  lab = cv2.merge(lab_planes)

  bgr = cv2.cvtColor(lab, cv2.COLOR_LAB2BGR)


  size = (256, 256)
  bgr = resize(bgr, size, mode='constant', anti_aliasing=True)
  return bgr

def thinning(bw):
  square_parameters = [2]
  bw_thinned = 0
  for i in range(len(square_parameters)):
    bw_thinned = binary_erosion(bw, footprint=np.ones((square_parameters[i], square_parameters[i])))
  return bw_thinned

def make_resized_bw_image (readed_image):
  print("before resized image")
  size = (256, 256)
  image = resize(readed_image, size, mode='constant', anti_aliasing=True)
  print("after resized image")

  save_image(image, 'nearest.jpg')

  # Синий канал
  print("before blue")
  gray_image = image[:,:,2] 
  save_image(gray_image, 'nearest_gray.jpg', True)
  print("after blue")

  print("before CLAHE")
  print("type(readed_image) = " + str(type(readed_image)))
  print("readed_image.shape = "+ str(readed_image.shape))
  bgr = CLAHE_to_image(readed_image)
  save_image(bgr, 'clahe.jpg', True)
  print("afte Clahe")


  threshold = 0.25 # 0.
  bw = bgr < threshold
  bw = bw[:,:,2]
  save_image(bw, 'threshold.jpg', True)
  return bw


def lesion_subtraction (full_image, lesion_image):
  full_image = full_image.tolist()
  lesion_image = lesion_image.tolist()
  for i in range(len(full_image)):
    for j in range(len(full_image[i])):
      if full_image[i][j] == True and lesion_image[i][j] == True:
        full_image[i][j] = False
  return np.array(full_image)

def count_regions (colored):
  # Подсчет участков
  label_img = label(colored)
  regions = regionprops(label_img)
  region_numbers = len(regions)

  return region_numbers

def count_radius(image_bw):
    """
    Функция для поиска и рисования круга на изображении.
    :param img_path: путь к изображению.
    """
    print()
    print("Circle Start")
    # Чтение изображения и преобразование его в оттенки серого.

    image = np.empty(image_bw.shape, dtype=np.uint8)
    for i in range(image_bw.shape[0]):
      for j in range (image_bw.shape[1]):
        if image_bw[i][j] == False:
          image[i][j] = 0
        else:
          image[i][j] = 255

    print("type(img)", type(image))
    print("img.dtype", image.dtype)
    print("img.shape", image.shape)
    
    # Поиск контуров на изображении.
    contours, _ = cv2.findContours(image, cv2.RETR_LIST, cv2.CHAIN_APPROX_SIMPLE)
    
    # Нахождение круга, описывающего контур.

    # get max circle 
    max_radius = 0
    max_countour = None
    max_center = None

    for i in range (len(contours)):
      (x,y), radius = cv2.minEnclosingCircle(contours[i])
      center = (int(x),int(y))
      radius = int(radius)
      if max_countour == None or max_radius < radius:
        max_radius = radius
        max_countour = contours[i]
        max_center = center

    max_area = cv2.contourArea(max_countour)
    
    print('Radius = ' + str(max_radius))
    print('Area = ' + str(max_area))
    
    return max_radius, max_area

def thinning_and_count_regions(bw):

  bw_thinned = thinning(bw)
  save_image(bw_thinned, 'bw_thinned.jpg', True)

  return count_regions(bw_thinned)

def segment_lesion_bw(extension, image_size, image_data):
  print()
  print("type(image_data)", type(image_data))
  print("len(image_data)", len(image_data))
  channel = grpc.insecure_channel('[::]:40043')
  stub = lesionSegmentation_pb2_grpc.LesionSegmentationStub(channel)
  request = lesionSegmentation_pb2.SegmentLesionRequest(extenstion = extension, height = image_size[0], width = image_size[1], image = image_data)
  response = stub.Segment(request)
  print("Before open bytes of segmented image")
  
  deserialized_bytes = np.frombuffer(response.segmentedImage, dtype=np.int64)
  segmented_image = np.reshape(deserialized_bytes, newshape=(256, 256))


  print("After open bytes of segmented image")
  
  print(type(segmented_image))
  return segmented_image

def lesion_proportion(image):
  unique = np.unique(image, return_counts=True)
  print("Unique values in image array", unique)
  true_quantity = unique[1][0]
  false_quantity = unique[1][1]
  return true_quantity/(false_quantity + true_quantity)

def count_abs_parameters(pores_density, pores_quantity, lesion_bw, lesion_proportion, area, radius):
  print("count_abs_parameters begin")
  abs_full_area =  pores_quantity / pores_density / (1 - lesion_proportion)
  abs_lesion_area = abs_full_area * lesion_proportion
  print()
  print("lesion_proportion = ", lesion_proportion)
  print()
  height = lesion_bw.shape[0]
  width = lesion_bw.shape[1]
  abs_pixel_area = abs_full_area / (height * width)
  abs_pixel_side = math.sqrt(abs_pixel_area)

  abs_radius = radius * abs_pixel_side
  return abs_lesion_area, abs_radius



def count_lesion_parameters(lesion_bw, pores_quantity, density):
  lesion_fraction = lesion_proportion(lesion_bw)
  outside_area = density * pores_quantity
  lesion_area = outside_area/(1 - lesion_fraction) * lesion_fraction
  
  radius, area = count_radius(lesion_bw)

  abs_lesion_area, abs_radius = count_abs_parameters(density, pores_quantity, lesion_bw, lesion_fraction, area, radius)
  abs_diameter = abs_radius * 2
  return abs_lesion_area, abs_diameter

def count_pores(extension, image_size, data, density):
  print()
  print("count_pores start")
  print(extension)
  print(image_size)
  print(type(data))
  print("len(data)", len(data))


  readed_image = Image.open(io.BytesIO(data))
  print(type(readed_image))
  readed_image = np.asarray(readed_image)
  bw = make_resized_bw_image(readed_image)
  save_image(bw, 'bw.jpg', True)
  print("after save_image")
  
  lesion_bw = segment_lesion_bw(extension, image_size, data)

  lesion_bw_bool = np.empty(lesion_bw.shape, dtype=bool)
  for i in range(256):
    for j in range(256):
      if lesion_bw[i][j] == 0:
        lesion_bw_bool[i][j] = False
      else:
        lesion_bw_bool[i][j] = True

  save_image(lesion_bw, 'lesion_3_channels_bw.jpg', True)
  print("after saving, lesion_bw " + str(type(lesion_bw)))
  
  save_image(lesion_bw_bool, 'lesion_bw_bool.jpg', True)

  # substract lesion
  bw_without_lesion = lesion_subtraction(bw, lesion_bw_bool)
  save_image(bw_without_lesion, 'bw_without_lesion.jpg', True)

  # count pores quantity
  pores_quantity = thinning_and_count_regions(bw_without_lesion)
  print('Кол-во регионов = ' + str(pores_quantity))
  area, diameter = count_lesion_parameters(lesion_bw, pores_quantity, density=density)
  print("count_pores end")
  return area, diameter

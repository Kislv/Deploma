import torch
import torch.nn as nn
import numpy as np
import matplotlib.pyplot as plt
import os

from skimage.io import imread
from torch.utils.data import DataLoader
from skimage.transform import resize
from PIL import Image
import io

class UNet(nn.Module):
    def __init__(self):
        super(UNet, self).__init__()

        self.enc_conv0 = nn.Sequential(
            nn.Conv2d(in_channels=3, out_channels=64, kernel_size=3, padding=1),
            nn.ReLU(),
            nn.Conv2d(in_channels=64, out_channels=64, kernel_size=3, padding=1),
            nn.ReLU()
        )

        self.pool0 = nn.MaxPool2d(kernel_size=2) # 256 -> 128

        self.enc_conv1 = nn.Sequential(
            nn.Conv2d(in_channels=64, out_channels=128, kernel_size=3, padding=1),
            nn.ReLU(),
            nn.Conv2d(in_channels=128, out_channels=128, kernel_size=3, padding=1),
            nn.ReLU()
        )
        self.pool1 = nn.MaxPool2d(kernel_size=2) # 128 -> 64

        self.enc_conv2 = nn.Sequential(
            nn.Conv2d(in_channels=128, out_channels=256, kernel_size=3, padding=1),
            nn.ReLU(),
            nn.Conv2d(in_channels=256, out_channels=256, kernel_size=3, padding=1),
            nn.ReLU()
        )
        self.pool2 = nn.MaxPool2d(kernel_size=2) # 64 -> 32

        self.enc_conv3 = nn.Sequential(
            nn.Conv2d(in_channels=256, out_channels=512, kernel_size=3, padding=1),
            nn.ReLU(),
            nn.Conv2d(in_channels=512, out_channels=512, kernel_size=3, padding=1),
            nn.ReLU()
        )
        self.pool3 = nn.MaxPool2d(kernel_size=2) # 32 -> 16

        # bottleneck
        self.bottleneck_conv = nn.Sequential(
            nn.Conv2d(in_channels=512, out_channels=1024, kernel_size=3, padding=1),
            nn.ReLU(),
            nn.Conv2d(in_channels=1024, out_channels=512, kernel_size=3, padding=1),
            nn.ReLU()
        )

        # decoder (upsampling)
        self.upsample0 = nn.Upsample(scale_factor=2, mode='nearest') # 16 -> 32
        self.dec_conv0 = nn.Sequential(
            nn.Conv2d(in_channels=1024, out_channels=512, kernel_size=3, padding=1),
            nn.ReLU(),
            nn.Conv2d(in_channels=512, out_channels=256, kernel_size=3, padding=1),
            nn.ReLU()
        )

        self.upsample1 = nn.Upsample(scale_factor=2, mode='nearest') # 32 -> 64
        self.dec_conv1 = nn.Sequential(
            nn.Conv2d(in_channels=512, out_channels=256, kernel_size=3, padding=1),
            nn.ReLU(),
            nn.Conv2d(in_channels=256, out_channels=128, kernel_size=3, padding=1)
        )

        self.upsample2 = nn.Upsample(scale_factor=2, mode='nearest')  # 64 -> 128
        self.dec_conv2 = nn.Sequential(
            nn.Conv2d(in_channels=256, out_channels=128, kernel_size=3, padding=1),
            nn.ReLU(),
            nn.Conv2d(in_channels=128, out_channels=64, kernel_size=3, padding=1),
            nn.ReLU()
        )
        
        self.upsample3 = nn.Upsample(scale_factor=2, mode='nearest')  # 128 -> 256
        self.dec_conv3 = nn.Sequential(
            nn.Conv2d(in_channels=128, out_channels=64, kernel_size=3, padding=1),
            nn.ReLU(),
            nn.Conv2d(in_channels=64, out_channels=64, kernel_size=3, padding=1),
            nn.ReLU(),
            nn.Conv2d(in_channels=64, out_channels=1, kernel_size=1)
        )

    def forward(self, x):
        # encoder
        e0 = self.enc_conv0(x)
        e0_pool = self.pool0(e0)

        e1 = self.enc_conv1(e0_pool)
        e1_pool = self.pool1(e1)

        e2 = self.enc_conv2(e1_pool)
        e2_pool = self.pool2(e2)

        e3 = self.enc_conv3(e2_pool)
        e3_pool = self.pool3(e3)

        # bottleneck
        b = self.bottleneck_conv(e3_pool)

        # decoder
        d0 = self.dec_conv0(torch.cat((self.upsample0(b), e3), dim=1))
        d1 = self.dec_conv1(torch.cat((self.upsample1(d0), e2), dim=1))
        d2 = self.dec_conv2(torch.cat((self.upsample2(d1), e1), dim=1))
        d3 = self.dec_conv3(torch.cat((self.upsample3(d2), e0), dim=1))

        return d3




def predict(model, data):
    model.eval()  # testing mode
    Y_pred = [model(X_batch) for X_batch in data]
    return np.array(Y_pred)


def show_before_after (device, model, data):  
  model.eval()  # testing mode
  with torch.no_grad():
    for X_val, Y_val in data:
      
      # data to device
      X_val = X_val.to(device)
      Y_val = Y_val.to(device)
      
      Y_pred = model(X_val)  
      
      Y_hat = Y_pred # model(X_val)
      plt.imshow(np.rollaxis(X_val[0].to('cpu').numpy(), 0, 3))
      plt.axis('off')
      plt.savefig("Images/before.png", bbox_inches='tight', pad_inches = 0)
    #   plt.show()

      Y_hat = torch.where(torch.sigmoid(Y_hat)>0.5,1,0)
      plt.imshow(Y_hat[0, 0].to('cpu').numpy(), cmap='gray')
      plt.savefig("Images/after.png", bbox_inches='tight', pad_inches = 0)
    #   plt.show()

      segmented_lesion = Y_hat[0, 0].to('cpu').numpy()
      print("type(segmented_lesion)", type(segmented_lesion))
      print("type(segmented_lesion[0][0])", type(segmented_lesion[0][0]))
      print("Before shape")
      print("segment_lesion.shape", segmented_lesion.shape)
      segmented_lesion_bytes = segmented_lesion.tobytes()
      print()
      print("type(segmented_lesion_bytes)", type(segmented_lesion_bytes))
      print("len(segmented_lesion_bytes)", len(segmented_lesion_bytes))
      return segmented_lesion_bytes
    
def segment_lesion(extension, image_size, image_data):
  SEED = 42
  ROOT = '/home/viktor/Projects/Univer/Sem8/Deploma/Deploma'
  LESION_SEGMENTATION_MODEL_PATH = os.path.join(ROOT, 'ML/LesionSegmentation/Models/Unet2.pt')

  print()
  print("segment_lesion")
  print(extension)
  print(image_size)
  print("len(image_data)", len(image_data))


  readed_image = Image.open(io.BytesIO(image_data))
  print(type(readed_image))
  readed_image = np.asarray(readed_image)
  print("readed_image " + str(type(readed_image)))
  print("readed_image.shape" + str(readed_image.shape))

  device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')
  print(device)

  m_state_dict = torch.load(LESION_SEGMENTATION_MODEL_PATH, map_location=(device))
  model = UNet().to(device)
  model.load_state_dict(m_state_dict)
  model.eval()

  images = []
  lesions = []
  for i in range(6):
    print("type(readed_image)" + str(type(readed_image)))
    print("readed_image.shape", readed_image.shape)
    images.append(readed_image)
    lesions.append(readed_image)

  size = (256, 256)
  X = [resize(x, size, mode='constant', anti_aliasing=True,) for x in images]
  Y = [resize(y, size, mode='constant', anti_aliasing=False) > 0.5 for y in lesions]

  X = np.array(X, np.float32)
  Y = np.array(Y, np.float32)
  print(f'Loaded {len(X)} images')

  len(lesions)

  np.random.seed(SEED)

  ix = np.random.choice(len(X), len(X), False)
  tr, val, ts = np.split(ix, [2, 4])

  print(len(tr), len(val), len(ts))

  batch_size = 10
  data_val = DataLoader(list(zip(np.rollaxis(X[val], 3, 1), Y[val, np.newaxis])),
                        batch_size=batch_size, shuffle=True)


  return show_before_after(device, model, data_val)

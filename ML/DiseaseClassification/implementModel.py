import os
import io
from PIL import Image
from matplotlib import pyplot as plt
import numpy as np
from torchvision import datasets, transforms
import torch
import torchvision.models as models



def printHello():
    print("hello from other file!")

def classificateImage(extension, size,  data):
    print('!!!!!!!!')
    print(size)
    print(type(data))
    ROOT = '/home/viktor/Projects/Univer/Sem8/Deploma/Deploma'
    DISEASE_SEGMENTATION = 'ML/DiseaseClassification'
    IMAGES = 'Images'
    MODELS = 'Models'
    IMAGE_CLASSIFICATION_PATH = os.path.join(IMAGES, 'mole_pores2.jpg')

    print(f'len(data) = {len(data)}')
    print(f'size = {size}')
    print(f'size[0] * size[1] = {size[0] * size[1]}')
    
    readed_image = Image.open(io.BytesIO(data))
    readed_image_tensor = transforms.ToTensor()(readed_image)
    print(type(readed_image))
    print(readed_image)
    plt.imshow(np.rollaxis(readed_image_tensor.to('cpu').numpy(), 0, 3))
    plt.axis('off')
    plt.savefig("FROM_SERVER.png", bbox_inches='tight', pad_inches = 0)
    
    readed_image_tensor = readed_image_tensor.unsqueeze(0)

    device = torch.device('cuda' if torch.cuda.is_available() else 'cpu')

    DISEASE_CLASSIFICATION_MODEL_PATH = os.path.join(ROOT, DISEASE_SEGMENTATION, MODELS, 'swin_v2.pt')
    m_state_dict = torch.load(DISEASE_CLASSIFICATION_MODEL_PATH, map_location=(device))

    del m_state_dict['fc.bias']
    del m_state_dict['fc.weight']

    pretrained_swin = models.swin_v2_t(weights=models.Swin_V2_T_Weights.DEFAULT)


    loaded_swin = (pretrained_swin).to(device)
    loaded_swin.load_state_dict(m_state_dict)


    classes = ['Hair Loss Photos Alopecia and other Hair Diseases',
    'Bullous Disease Photos',
    'Scabies Lyme Disease and other Infestations and Bites',
    'Psoriasis pictures Lichen Planus and related diseases',
    'Acne and Rosacea Photos',
    'Systemic Disease',
    'Tinea Ringworm Candidiasis and other Fungal Infections',
    'Actinic Keratosis Basal Cell Carcinoma and other Malignant Lesions',
    'Eczema Photos',
    'Warts Molluscum and other Viral Infections',
    'Urticaria Hives',
    'Seborrheic Keratoses and other Benign Tumors',
    'Lupus and other Connective Tissue diseases',
    'Herpes HPV and other STDs Photos',
    'Vasculitis Photos',
    'Cellulitis Impetigo and other Bacterial Infections',
    'Nail Fungus and other Nail Disease',
    'Poison Ivy Photos and other Contact Dermatitis',
    'Melanoma Skin Cancer Nevi and Moles',
    'Vascular Tumors',
    'Exanthems and Drug Eruptions',
    'Atopic Dermatitis Photos',
    'Light Diseases and Disorders of Pigmentation']

    outputs = loaded_swin(readed_image_tensor.to(device))
    predicted = torch.argmax(outputs.data,-1)    
    print("predicted[0]", predicted[0])
    print(classes[predicted.item()])   
    return classes[predicted.item()]


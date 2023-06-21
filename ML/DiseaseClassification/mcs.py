from concurrent import futures
import sys
import grpc

import implementModel 
sys.path.append( '../../GRPC/DiseaseClassification/ML' )
import image_pb2
import image_pb2_grpc



class ClassificateImageServicer(image_pb2_grpc.ClassificateImageServicer):

  def Classificate(self, request, context):
    print(type(request.image))
    predict = implementModel.classificateImage(request.extenstion, (request.width, request.height), request.image)
    print(request.width)
    print(request.height)
    return image_pb2.ImageClassificationResponse(DiseaseName = predict)
  
def serve():
  server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
  image_pb2_grpc.add_ClassificateImageServicer_to_server(
      ClassificateImageServicer(), server)
  server.add_insecure_port('[::]:40041')
  server.start()
  print("Disease classification mcs start!")
  server.wait_for_termination()

serve()

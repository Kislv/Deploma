from concurrent import futures
import sys
import grpc

import implementModel 
sys.path.append( '../../GRPC/LesionSegmentation/ML' )
import lesionSegmentation_pb2
import lesionSegmentation_pb2_grpc


class LesionSegmentationServicer(lesionSegmentation_pb2_grpc.LesionSegmentationServicer):

  def Segment(self, request, context):
    print(type(request.image))
    print(request.width)
    print(request.height)
    segmented_lesion = implementModel.segment_lesion(request.extenstion, (request.width, request.height), request.image)
    return lesionSegmentation_pb2.SegmentLesionResponse(segmentedImage = segmented_lesion)
  
def serve():
  server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
  lesionSegmentation_pb2_grpc.add_LesionSegmentationServicer_to_server(
      LesionSegmentationServicer(), server)
  server.add_insecure_port('[::]:40043')
  server.start()
  print("Lesion segmentation mcs start!")
  server.wait_for_termination()

serve()

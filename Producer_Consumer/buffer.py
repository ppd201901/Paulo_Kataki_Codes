import xmlrpclib
from SimpleXMLRPCServer import SimpleXMLRPCServer
from multiprocessing.pool import ThreadPool
from multiprocessing import Lock


BUFFER_SIZE = 10
indexBuffer = -1

mutex = Lock()

buffer = range(BUFFER_SIZE)

server = SimpleXMLRPCServer(("localhost", 8080))


def Store(data):
	global mutex,indexBuffer
	mutex.acquire()
	if (indexBuffer == BUFFER_SIZE - 1):
		mutex.release()
		return "Full buffer"
	else:
		indexBuffer += 1
		buffer[indexBuffer] = data
		print ("Stored sentence: \"" + data + "\" sucessfully")
		mutex.release()
		return "Stored sentence: \"" + data + "\" sucessfully"


def Store_RPC(data):
	async_result = pool.apply_async(Store, (data,))
	return async_result.get() 



def Consume():
	global mutex,indexBuffer
	mutex.acquire()
	if (indexBuffer == -1):
		mutex.release()
		return "Empty buffer"
	else :
		indexBuffer -= 1
		mutex.release()
		print "Sentence consumed: " + buffer[indexBuffer + 1]
		return (buffer[indexBuffer+1] )


def Consume_RPC():
	async_result = pool.apply_async(Consume)
	return async_result.get()


pool = ThreadPool(processes=10)


server.register_function(Store_RPC, "produce")
server.register_function(Consume_RPC, "consume")
print "Server up and running"


server.serve_forever()

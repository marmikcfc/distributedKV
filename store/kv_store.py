import threading
lock = threading.Lock()

shared_dict = {}

class KVStore():

  def __init__(self):
    self.store = {}

  def get(self, key):
    lock.acquire()
    try:
      ret = self.store.get(key,"Key does not exist")
    finally:
      lock.release()
    return ret

  def put(self, key, value):
    lock.acquire()
    try:
      self.store[key] = value
    finally:
      lock.release()
      
    return True

  def delete(self,key):
    try:
        del self.store[key]
    except KeyError:
        result_msg = 'ERROR: Key [{}] not found and could not be deleted'.format(key)
    else:
        result_msg = "Key [{}] deleted".format(key)

    return result_msg

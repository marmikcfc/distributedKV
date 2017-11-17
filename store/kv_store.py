class KVStore():

  def __init__(self):
    self.store = {}

  def get(self, key):
    return self.store.get(key, key + " not found")

  def put(self, key, value):
    self.store[key] = value
    return True

  def delete(self,key):
    try:
        del self.store[key]
    except KeyError:
        result_msg = 'ERROR: Key [{}] not found and could not be deleted'.format(key)
    else:
        result_msg = "Key [{}] deleted".format(key)

    return result_msg

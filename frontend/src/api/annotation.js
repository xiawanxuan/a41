import request from './request'

export function createAnnotation(data) {
  return request({
    url: '/annotations',
    method: 'post',
    data
  })
}

export function getAnnotation(id) {
  return request({
    url: `/annotations/${id}`,
    method: 'get'
  })
}

export function getAnnotationsByText(textId, userId) {
  return request({
    url: `/annotations/text/${textId}`,
    method: 'get',
    params: { user_id: userId }
  })
}

export function updateAnnotation(id, data) {
  return request({
    url: `/annotations/${id}`,
    method: 'put',
    data
  })
}

export function deleteAnnotation(id) {
  return request({
    url: `/annotations/${id}`,
    method: 'delete'
  })
}

export function batchSubmitAnnotations(data) {
  return request({
    url: '/annotations/batch',
    method: 'post',
    data
  })
}

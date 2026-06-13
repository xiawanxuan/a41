import request from './request'

export function getLayoutByText(textId) {
  return request({
    url: `/layouts/text/${textId}`,
    method: 'get'
  })
}

export function saveLayout(textId, data) {
  return request({
    url: `/layouts/text/${textId}`,
    method: 'put',
    data
  })
}

export function deleteLayout(id) {
  return request({
    url: `/layouts/${id}`,
    method: 'delete'
  })
}

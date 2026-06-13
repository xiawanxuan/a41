import request from './request'

export function createText(data) {
  return request({
    url: '/texts',
    method: 'post',
    data
  })
}

export function getText(id) {
  return request({
    url: `/texts/${id}`,
    method: 'get'
  })
}

export function listTexts(params) {
  return request({
    url: '/texts',
    method: 'get',
    params
  })
}

export function updateText(id, data) {
  return request({
    url: `/texts/${id}`,
    method: 'put',
    data
  })
}

export function deleteText(id) {
  return request({
    url: `/texts/${id}`,
    method: 'delete'
  })
}

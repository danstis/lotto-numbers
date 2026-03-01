const { test } = require('node:test')
const assert = require('node:assert/strict')
const { getBallClass } = require('./app.js')

test('getBallClass returns ball-blue for numbers 1-9', () => {
  assert.equal(getBallClass(1), 'ball-blue')
  assert.equal(getBallClass(5), 'ball-blue')
  assert.equal(getBallClass(9), 'ball-blue')
})

test('getBallClass returns ball-orange for numbers 10-19', () => {
  assert.equal(getBallClass(10), 'ball-orange')
  assert.equal(getBallClass(15), 'ball-orange')
  assert.equal(getBallClass(19), 'ball-orange')
})

test('getBallClass returns ball-green for numbers 20-29', () => {
  assert.equal(getBallClass(20), 'ball-green')
  assert.equal(getBallClass(25), 'ball-green')
  assert.equal(getBallClass(29), 'ball-green')
})

test('getBallClass returns ball-red for numbers 30-39', () => {
  assert.equal(getBallClass(30), 'ball-red')
  assert.equal(getBallClass(35), 'ball-red')
  assert.equal(getBallClass(39), 'ball-red')
})

test('getBallClass returns ball-purple for number 40', () => {
  assert.equal(getBallClass(40), 'ball-purple')
})

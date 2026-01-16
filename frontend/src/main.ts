import { mount } from 'svelte'
import App from './App.svelte'
import Settings from './Settings.svelte'

// 根据 URL hash 决定加载哪个页面
const hash = window.location.hash

const target = document.getElementById('app')!

let app: any

if (hash === '#/settings') {
  app = mount(Settings, { target })
} else {
  app = mount(App, { target })
}

export default app

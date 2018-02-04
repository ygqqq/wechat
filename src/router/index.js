import Vue from 'vue'
import Router from 'vue-router'
import Login from '@/pages/Login'
import Register from '@/pages/Register'
import Home from '@/pages/Home'
import Addres from '@/pages/Addres'
import Message from '@/pages/Message'
import Chat from '@/pages/Chat'

Vue.use(Router)

export default new Router({
  mode:'history',
  routes: [
    {
      path: '/',
      redirect: '/message'
    },
    {
      path: '/login',
      name: 'login',
      component: Login
    },
    {
      path: '/register',
      name: 'register',
      component: Register
    },
    {
      path: '/home',
      name: 'home',
      component: Home
    },
    {
      path: '/addres',
      name: 'addres',
      component: Addres
    },
    {
      path: '/message',
      name: 'message',
      component: Message
    },
    {
      path: '/chat/:username',
      name: 'chat',
      component: Chat
    }
  ]
})

// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import MintUI from 'mint-ui'
import 'mint-ui/lib/style.css'
Vue.use(MintUI)

Vue.config.productionTip = false


/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: { App },
  template: '<App/>'
})

// OnlineRemind	= 1	//上线提醒
// OfflineRemind   = 2 //下线提醒 
// AddFriendReq	= 3 //添加好友请求
// AgreeAdd		= 4 //同意好友请求
// DisAgreeAdd 	= 5 //拒绝好友请求
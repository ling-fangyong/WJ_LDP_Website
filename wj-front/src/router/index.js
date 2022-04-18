import Vue from 'vue'
import Router from 'vue-router'
import Register from '../views/Register.vue'
import Login from '../views/Login.vue'
import TopBar from '../views/TopBar.vue'
import Home from '../views/Home.vue'
import Display from '../views/display.vue'

Vue.use(Router)

export default new Router({
  mode:'history',
  routes: [
    {
      path: '/',
      name: 'TopBar',
      component: TopBar,
      children:[
        {
          path:'/',
          name:'Home',
          component: Home
        },
        {
          path:'/Home',
          name:'Home',
          component: Home
        },
        {
          path:'/Login',
          name:'Login',
          component:Login
        },
        {
          path:'/Register',
          name:'Register',
          component:Register
        },
    ]
    },
    {
      path:'/display/:id',
      name:'Display',
      component:Display,
    }
  ]
})

import { createRouter, createWebHistory } from 'vue-router'
import Login from '../views/auth/Login'
import AdminApp from '../layouts/AdminApp.vue'
import AdminBaseLayout from '../layouts/AdminBaseLayout';
import Dashboard from '../views/admin/dashboard/Dashboard';
import Workspaces from '../views/admin/workSpace/WorkSpace';
import UserBaseLayout from '../layouts/UserBaseLayout';
import UserDashboard from '../views/user/dashboard/UserDashboard'
import WorkspacesDetail from '../views/admin/workSpace/WorkSpaceDetail'
import VirtualMachineDetail from "../views/admin/VirtualMachine/VirtualMachineDetail";
import VirtualMachine from "../views/admin/VirtualMachine/VirtualMachine";

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: Login
  },
  {
    path: '/adminApp',
    name: 'AdminApp',
    component: AdminApp
  },
  {
    path: '/',
    name: 'home',
    component: AdminBaseLayout,
    meta: { icon: 'home' },
    redirect: '/dashboard',
    children: [
      {
        path: '/dashboard',
        name: 'Dashboard',
        component: Dashboard
      },
      {
        path: '/workspaces',
        name: 'Workspaces',
        component: Workspaces
      },
      {
        path: '/workspacesDetail/',
        name: 'WorkspacesDetail',
        component: WorkspacesDetail,
        props: true
      },
      {
        path: '/virtualmachine',
        name: 'Virtualmachine',
        component: VirtualMachine
      },
      {
        path: '/vmdetail/',
        name: 'VirtualMachineDetail',
        component: VirtualMachineDetail,
        props: true
      }
    ]
  },
  {
    path: '/user',
    name: 'User',
    component: UserBaseLayout,
    meta: { icon: 'home' },
    redirect: '/userDashboard',
    children: [
      {
        path: '/userDashboard',
        name: 'UserDashboard',
        component: UserDashboard
      }
    ]
  }
]

const index = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes
})

export default index

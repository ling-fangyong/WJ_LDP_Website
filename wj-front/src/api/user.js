import axios from 'axios';

const register = form => axios.post('/api/user/register',form).then(res => res.data);

const login = form => axios.post('/api/user/login', form).then(res => res.data);

const logout = () => axios.delete('/api/user/logout').then(res => res.data);

export {
    register,
    login,
    logout,
};
import axios from 'axios';

const UpdateQuestionaire = form => axios.post('/api/design/UpdateQuestionaire',form).then(res => res.data);

const UpdateQuestion = form => axios.post('/api/design/UpdateQuestion', form).then(res => res.data);

const DeleteQuestionaire = form => axios.delete('/api/design/DeleteQuestionaire',form).then(res => res.data);

const ShowQuestionaires = () => axios.get('/api/show/ShowQuestionaires').then(res => res.data);

const ShowQuestions = form => axios.post('/api/show/ShowQuestions',form).then(res => res.data)

export {
    UpdateQuestionaire,
    UpdateQuestion,
    DeleteQuestionaire,
    ShowQuestionaires,
    ShowQuestions
};
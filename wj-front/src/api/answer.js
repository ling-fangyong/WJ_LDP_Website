import axios from "axios";

const GetQuestionaire = form => axios.post('/api/answer/GetQuestionaire',form).then(res => res.data);

const SubmitQues =form => axios.post('/api/answer/SubmitQues',form).then(res => res.data);

const AnalysisData = form => axios.post('/api/analysis/AnalysisData',form).then(res => res.data);

export {
    GetQuestionaire,
    SubmitQues,
    AnalysisData,
}
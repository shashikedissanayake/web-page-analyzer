import { useCallback, useState } from "react";
import axios from 'axios';
import { Response } from "./Response";

export const WebForm = () => {
    const [text, setText] = useState('');
    const [data, setData] = useState(null);
    const [error, setError] = useState(null);

    const handleUrlChanges = (e) => {
        setText(e.target.value);
    }

    const request = useCallback(async (url) => {
        setData(null);
        setError(null);
        try {
            const response = await axios.post(`${process.env.REACT_APP_SERVER_URL}/analyze`, { url });
            setData(response?.data?.data);
        } catch (error) {
            setError(error?.response?.data?.message || 'Failed to access given Url');
        }

    }, [setData, setError]);

    const handleSubmit = (e) => {
        e.preventDefault();
        request(text);
        setText('');
    }

    return (
        <div>
            <form onSubmit={handleSubmit}>
                <div>
                    <input
                        type="text"
                        value={text}
                        onChange={handleUrlChanges}
                    />
                    <button>submit</button>

                </div>
            </form>
            <div>
                {(data || error) && <Response data={data} error={error} />}
            </div>
        </div>
    );
}
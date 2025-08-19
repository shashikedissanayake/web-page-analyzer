export const Response = ({
    error,
    data,
}) => {
    return (
        <div>
            {
                error &&
                <div>
                    <h2>Failed to analyze URL</h2>
                    <p>{error}</p>
                </div>
            }
            {
                data && <div>
                    <h1>Successfully analyzed URL</h1>
                    <h3>URL: {data.url}</h3>
                    <h3>Html version: {data.htmlVersion}</h3>
                    <h3>Page title: {data.pageTitle}</h3>
                    <h3>Is a login form: {data.isLoginForm ? 'yes' : 'no'}</h3>
                    <h3>Internal links: Found internal links {data.internalLinkCount} and {data.internalInaccessibleLinkCount} links are inaccessible</h3>
                    <h3>External links: Found external links {data.externalLinkCount} and {data.externalInaccessibleLinkCount} links are inaccessible</h3>
                    <ul>
                        {
                            Object.entries(data?.headerTagCount || {})?.map(
                                (
                                    { tag, count }) => {
                                    return (
                                        <li>
                                            <h3>Found {count} tags of {tag}</h3>
                                        </li>
                                    )
                                }
                            )
                        }
                    </ul>
                </div>
            }
        </div>
    );
}

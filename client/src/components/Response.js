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
                    <h3>Internal link count: {data.internalLinkCount}</h3>
                    <h3>Inaccessible Internal link count: {data.internalInaccessibleLinkCount}</h3>
                    <h3>External link count: {data.externalLinkCount}</h3>
                    <h3>Inaccessible External link count: {data.externalInaccessibleLinkCount}</h3>
                    {
                        Object.entries(data?.headerTagCount || {}).length > 0 ?
                            <>
                                <h3>Found following headers in the page</h3>
                                <ul>
                                    {
                                        Object.entries(data?.headerTagCount || {})?.map(
                                            ([tag, count]) => (
                                                <li key={tag}>
                                                    <h3>Found {count} of {tag} tags</h3>
                                                </li>
                                            )
                                        )
                                    }
                                </ul>
                            </>
                            : <h3>No header tags found</h3>
                    }
                </div>
            }
        </div>
    );
}

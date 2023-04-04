const allowed = [
  'acct:john@velvetcache.org',
  'acct:jmhobbs@velvetcache.org',
];

export function onRequest(context) {
  const url = new URL(context.request.url);
  console.log(url.searchParams.get('resource'));
  if(allowed.includes(url.searchParams.get('resource'))) {
    return Response.redirect('https://noc.social/.well-known/webfinger?resource=acct:jmhobbs@noc.social', 302);
  }
  return new Response('400 Bad Request', { status: 400, statusText: 'Bad Request' });
}

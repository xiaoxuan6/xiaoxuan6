import os
import re

import feedparser


def replace_chunk(content, marker, chunk, inline=False):
    r = re.compile(
        r"<!\-\- {} starts \-\->.*<!\-\- {} ends \-\->".format(marker, marker),
        re.DOTALL,
    )
    if not inline:
        chunk = "\n{}\n".format(chunk)
    chunk = "<!-- {} starts -->{}<!-- {} ends -->".format(marker, chunk, marker)
    return r.sub(chunk, content)


def sync():
    newsFeed = feedparser.parse('https://xiaoxuan6.github.io/rss.xml')
    entry = newsFeed.entries[:5]

    entries_md = []
    for detail in entry:
        published = f'{detail.published}'.split('T')[0]
        entries_md.append(f"<a href='{detail.link}' target='_blank'>{detail.title}</a> - {published}<br/>\n")

    filename = os.getcwd() + '/README.md'
    with open(filename, encoding='utf-8', mode='r') as f:
        rewritten = replace_chunk(f.read(), "blog", ''.join(entries_md))

        with open(filename, encoding='utf-8', mode='w') as f:
            f.write(rewritten)


if __name__ == '__main__':
    sync()

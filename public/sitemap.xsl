<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet version="2.0"
                xmlns:html="http://www.w3.org/TR/REC-html40"
                xmlns:sitemap="http://www.sitemaps.org/schemas/sitemap/0.9"
                xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
    <xsl:output method="html" version="1.0" encoding="UTF-8" indent="yes"/>
    <xsl:template match="/">
        <html xmlns="http://www.w3.org/1999/xhtml">
            <head>
                <title>XML Sitemap</title>
                <meta http-equiv="content-type" content="text/html; charset=utf-8" />
                <style type="text/css">
                    body { font-family: "Consolas", "Lucida Console", monospace; font-size: 11pt; color: black; }
                    a { color:#7F923F; }
                    h1 { font-size: 20pt; }
                    #intro { margin: 5px; padding: 5px 15px; color:gray; line-height: 1.0; }
                    #intro span { color: black; font-weight: bold; }
                    th { text-align:left; padding-right:15px; font-size: 12pt; border-bottom:1px black solid; }
                    td { text-align:center; }
                    td.url { text-align:left; }
                    tr.even { background-color: #D6E0B7; }
                    tr.even a { color: #59662C; }
                    #footer { padding:2px; margin:10px; font-size: 10pt; color:gray; }
                </style>
            </head>
            <body>
                <h1>XML Sitemap</h1>
                <div id="intro">
                    Вы наблюдаете слегка облагороженный сайтмап, пусть на голый XML роботы смотрят :)<br/>
                    Карта сайта содержит ссылок: <span><xsl:value-of select="count(sitemap:urlset/sitemap:url)"/></span>
                </div>
                <div id="content">
                    <table cellpadding="5">
                        <tr>
                            <th>URL</th>
                            <th>Последнее изменение</th>
                        </tr>
                        <xsl:for-each select="sitemap:urlset/sitemap:url">
                            <tr>
                                <xsl:if test="position() mod 2 != 1">
                                    <xsl:attribute  name="class">even</xsl:attribute>
                                </xsl:if>
                                <td class="url">
                                    <xsl:variable name="itemURL">
                                        <xsl:value-of select="sitemap:loc"/>
                                    </xsl:variable>
                                    <a href="{$itemURL}">
                                        <xsl:value-of select="sitemap:loc"/>
                                    </a>
                                </td>
                                <td>
                                    <xsl:value-of select="concat(substring(sitemap:lastmod,0,11),' ', substring(sitemap:lastmod,12,5))"/>
                                </td>
                            </tr>
                        </xsl:for-each>
                    </table>
                </div>
                <div id="footer">
                    Сгенерировано исключительно для красоты, хоть сюда всё равно никто смотреть и не будет.
                </div>
            </body>
        </html>
    </xsl:template>
</xsl:stylesheet>

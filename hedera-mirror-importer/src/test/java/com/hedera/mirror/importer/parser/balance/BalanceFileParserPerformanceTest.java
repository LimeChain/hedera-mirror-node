package com.hedera.mirror.importer.parser.balance;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import javax.annotation.Resource;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Timeout;
import org.junit.jupiter.api.io.TempDir;
import org.springframework.beans.factory.annotation.Value;
import org.testcontainers.shaded.org.apache.commons.io.FilenameUtils;

import com.hedera.mirror.importer.FileCopier;
import com.hedera.mirror.importer.IntegrationTest;
import com.hedera.mirror.importer.TestUtils;
import com.hedera.mirror.importer.domain.AccountBalanceFile;
import com.hedera.mirror.importer.domain.EntityId;
import com.hedera.mirror.importer.domain.StreamType;
import com.hedera.mirror.importer.repository.AccountBalanceFileRepository;
import com.hedera.mirror.importer.util.Utility;

@Tag("performance")
public class BalanceFileParserPerformanceTest extends IntegrationTest {

    private static final String DATA_SOURCE_FOLDER = "v1/performance";

    @TempDir
    static Path dataPath;

    @Value("classpath:data")
    Path testPath;

    @Resource
    private BalanceFileParser balanceFileParser;

    @Resource
    private BalanceParserProperties parserProperties;

    @Resource
    private AccountBalanceFileRepository accountBalanceFileRepository;

    private FileCopier fileCopier;

    private StreamType streamType;

    @BeforeEach
    void before() throws IOException {
        streamType = parserProperties.getStreamType();
        parserProperties.getMirrorProperties().setDataPath(dataPath);
        parserProperties.init();

        EntityId nodeAccountId = EntityId.of(TestUtils.toAccountId("0.0.3"));
        Files.walk(Path.of(testPath.toString(), streamType.getPath(), DATA_SOURCE_FOLDER))
                .filter(p -> p.toString().endsWith(".csv"))
                .forEach(p -> {
                    String filename = FilenameUtils.getName(p.toString());
                    AccountBalanceFile accountBalanceFile = AccountBalanceFile.builder()
                            .consensusTimestamp(Utility.getTimestampFromFilename(filename))
                            .count(0L)
                            .fileHash(filename)
                            .loadEnd(0L)
                            .loadStart(0L)
                            .name(filename)
                            .nodeAccountId(nodeAccountId)
                            .build();
                    accountBalanceFileRepository.save(accountBalanceFile);
                });
    }

    @Timeout(15)
    @Test
    void parseAndIngestMultipleBalanceCsvFiles() {
        parse("*.csv");
    }

    private void parse(String filePath) {
        fileCopier = FileCopier.create(testPath, dataPath)
                .from(streamType.getPath(), DATA_SOURCE_FOLDER)
                .filterFiles(filePath)
                .to(streamType.getPath(), streamType.getValid());
        fileCopier.copy();

        balanceFileParser.parse();
    }
}
